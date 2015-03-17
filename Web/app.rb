require 'rubygems'
require 'bundler/setup'
require 'bcrypt'

require './src/utils'
require './src/js'
require './src/db_user'
require './src/db_position'
require './src/db_archive'
require './src/db_analytics'

class BurgessApp < Sinatra::Base
	helpers Sinatra::JavaScripts

    configure do
        enable :sessions
        set :db_user, UserData.new
        set :db_position, PositionData.new
		set :db_archived, ArchiveData.new
		set :db_analytics, AnalyticsData.new
    end

    helpers do

        # Whether user is authenticated
        def authenticated?
            not session[:identity].nil?
        end

        # Get all errors and reset array
        def pop_errors
            tmp = session[:errors] || []
            session[:errors] = []
            return tmp
        end

        # Add error to array
        def push_error(error)
            (session[:errors] ||= []).push(error)
        end
    end

    get '/' do
        erb :home
    end

    get '/livefeed' do
		js :jcanvas, :knockout, 'map', 'knockout/livefeed'
        erb :livefeed
    end

	get '/livefeed/data' do
		session[:timelast] ||= Time.now - 5
		session[:timelast] = Time.now - 5 if session[:timelast] > Time.now - 10
		result = settings.db_archived.getPositionsSince(session[:timelast])
		session[:timelast] = Time.now
		return result.to_json
	end

    get '/playback' do
		js :datetime, :knockout, :jcanvas, :nvd3, 'map', 'timeselect', 'knockout/playback'
        erb :playback
    end

    get '/analytics' do
		js :nvd3, :knockout, :datetime, 'analytics/helpedTimeChart', 
            'analytics/peakChart', 'analytics/helpedCountChart', 'knockout/analytics'
        erb :analytics
    end

    get '/settings' do
		js :knockout, 'knockout/settings'
        erb :settings
    end

	### ANALYTICS ###

	post '/analytics/helpCount' do
		ti = Utils.StandardizeTime_s(params[:ti].to_i)
		tf = Utils.StandardizeTime_s(params[:tf].to_i)
		if authenticated?
			employees = settings.db_user.getEmployees(session[:identity].id)
			employees.map!{|e| e["_id"]}
            return settings.db_analytics.getEmployeeHelpCount(ti,tf,0,employees).to_json
		end
		return nil
	end

	post '/analytics/helpTime' do
		ti = Utils.StandardizeTime_s(params[:ti].to_i)
		tf = Utils.StandardizeTime_s(params[:tf].to_i)
		if authenticated?
			employees = settings.db_user.getEmployees(session[:identity].id)
            employees.map!{|e| e["_id"]}
			return settings.db_analytics.getEmployeeHelpTime(ti,tf,employees).to_json
		end
		return nil
	end

    post '/analytics/customersHourly' do
        ti = Utils.StandardizeTime_s(params[:ti].to_i)
        tf = Utils.StandardizeTime_s(params[:tf].to_i)
        if authenticated?
            return settings.db_archived.getVisitorsHourly(ti, tf, false).to_json
        end
        return nil
    end

    post '/analytics/employeesHourly' do
        ti = Utils.StandardizeTime_s(params[:ti].to_i)
        tf = Utils.StandardizeTime_s(params[:tf].to_i)
        if authenticated?
            return settings.db_archived.getVisitorsHourly(ti, tf, true).to_json
        end
        return nil
    end

	### EMPLOYEES ###

	get '/employees' do
		if authenticated?
			return settings.db_user.getEmployees(session[:identity].id).to_json
		end
		return nil
	end

	post '/employees' do
		if authenticated?
			settings.db_user.updateEmployees(session[:identity].id,Employee.fromArray(JSON.parse(params[:update])))
			settings.db_user.removeEmployees(session[:identity].id,Employee.fromArray(JSON.parse(params[:remove])))
		end
	end

	### MAP ###

    get '/map/size' do
        if authenticated?
            return session[:identity].getMapDetails.to_json
        end
		return nil
    end

    ### PLAY BACK ###

    post '/playback/date' do
        t = Utils.StandardizeTime_s(params[:t].to_i)
        timezone = params[:timezone].to_i
		return settings.db_archived.getPositionsOverDay(t, timezone).to_json
    end

    ### AUTHENTICATION ###

    get '/login' do
        erb :login
    end

    post '/login' do
        session[:identity] = settings.db_user.getUser(params['username'])
        if not session[:identity].nil? and session[:identity].validatePassword(params['password'])
            redirect to '/'
        else
            push_error('Invalid login')
            redirect to '/login'
        end
    end

    get '/logout' do
        session.delete(:identity)
        redirect to '/'
    end

    get '/signup' do
        erb :signup
    end

    post '/signup' do
        user = User.new.createUser(params['username'], params['password'], params['company'], params['storeID'])

        push_error("Username taken") if not settings.db_user.getUser(params['username']).nil?
        push_error("Passwords must match") if not user.validatePassword(params['re-password'])

        if not session[:errors] or session[:errors].empty?
            session[:identity] = user
            settings.db_user.storeUser(user)
            redirect to '/'
        else
            redirect to '/signup'
        end
    end
end
